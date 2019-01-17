package get

import (
	"fmt"
	"os"
	"strconv"

	"github.com/chengyumeng/khadijah/pkg/model"
	"github.com/chengyumeng/khadijah/pkg/utils/log"
	"github.com/olekukonko/tablewriter"
)

const pageSize int = 1024 * 1024 // 单页显示，不分页

type GetProxy struct {
	Option Option
}

func NewProxy(opt Option) GetProxy {
	return GetProxy{
		Option: opt,
	}
}

func (g *GetProxy) Get() {
	switch g.Option.Resource {
	case model.NamespaceType:
		g.getNamespace()
	case model.AppType:
		g.getApp()
	case model.DeploymentType:
		g.GetPod(model.DeploymentType)
	case model.StatefulsetType:
		g.GetPod(model.StatefulsetType)
	case model.DaemonsetType:
		g.GetPod(model.DaemonsetType)
	case model.CronjobType:
		g.GetPod(g.Option.Resource)
	case model.PodType:
		g.GetPod(model.DeploymentType)
		g.GetPod(model.StatefulsetType)
		g.GetPod(model.DaemonsetType)
		g.GetPod(model.CronjobType)
	case model.ServiceType:
		g.GetService()
	default:
	}
}

func (g *GetProxy) getNamespace() {
	data := model.GetNamespaceBody()
	fmt.Printf("Name: %s Email:%s\n\n", data.Data.Name, data.Data.Email)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Name", "User", "CreateTime", "UpdateTime"})

	for _, v := range data.Data.Namespaces {
		table.Append([]string{strconv.Itoa(int(v.Id)), v.Name, v.User, v.CreateTime.String(), v.UpdateTime.String()})
	}
	table.Render()
}

func (g *GetProxy) getApp() {
	nsIds := []int64{}
	ns := model.GetNamespaceBody()
	if g.Option.Namespace != "" {
		for _, n := range ns.Data.Namespaces {
			if n.Name == g.Option.Namespace {
				nsIds = append(nsIds, n.Id)
			}
		}
		if len(nsIds) == 0 {
			log.AppLogger.Warning("NS ERROR")
			return
		}
	} else {
		for _, n := range ns.Data.Namespaces {
			nsIds = append(nsIds, n.Id)
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Name", "Namespace", "User", "CreateTime"})
	for _, id := range nsIds {
		data := model.GetAppBody(id)
		if data == nil {
			continue
		}

		for _, v := range data.Data.Apps {
			table.Append([]string{strconv.Itoa(int(v.Id)), v.Name, v.Namespace, v.User, v.CreateTime.String()})
		}

	}
	table.Render()
}

func (g *GetProxy) GetPod(podType string) {
	nsIds := []int64{}
	ns := model.GetNamespaceBody()
	if g.Option.Namespace != "" {
		for _, n := range ns.Data.Namespaces {
			if n.Name == g.Option.Namespace {
				nsIds = append(nsIds, n.Id)
			}
		}
		if len(nsIds) == 0 {
			log.AppLogger.Warning("No NS")
		}
	} else {
		for _, n := range ns.Data.Namespaces {
			nsIds = append(nsIds, n.Id)
		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Name", "Type", "APP", "Namespace", "User", "CreateTime"})
	exist := false
	for _, nsId := range nsIds {
		if app := model.GetAppBody(nsId); app != nil {
			for _, a := range app.Data.Apps {
				if g.Option.App == "" || g.Option.App == a.Name {
					data := model.GetPodBody(a.Id, podType)
					for _, pod := range data.Data.Pods {
						exist = true
						table.Append([]string{strconv.Itoa(int(pod.Id)), pod.Name, podType, pod.App.Name, pod.App.NSMetaData.Name, pod.User, pod.CreateTime.String()})
					}
				}
			}
		}

	}
	if exist {
		table.Render()
	}
}

func (g *GetProxy) GetService() {
	nsl := []model.Namespace{}
	ns := model.GetNamespaceBody()
	if g.Option.Namespace != "" {
		for _, n := range ns.Data.Namespaces {
			if n.Name == g.Option.Namespace {
				nsl = append(nsl, n)
			}
		}
		if len(nsl) == 0 {
			log.AppLogger.Error("No NS")
		}
	} else {
		for _, n := range ns.Data.Namespaces {
			nsl = append(nsl, n)
		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Name", "Type", "APP", "Namespace", "User", "CreateTime"})
	exist := false
	for _, ns := range nsl {
		if app := model.GetAppBody(ns.Id); app != nil {
			for _, a := range app.Data.Apps {
				if g.Option.App == "" || g.Option.App == a.Name {
					data := model.GetServiceBody(a.Id)
					for _, svc := range data.Data.Services {
						exist = true
						table.Append([]string{strconv.Itoa(int(svc.Id)), svc.Name, model.ServiceType, a.Name, ns.Name, svc.User, svc.CreateTime.String()})
					}
				}

			}
		}
	}
	if exist {
		table.Render()
	}
}
