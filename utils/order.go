package utils

import "strings"

var AllowedOrders = map[string]bool{
	"created_at asc":  true,
	"created_at desc": true,
	"name asc":        true,
	"name desc":       true,
	"id asc":          true,
	"id desc":         true,
	// 可擴充其他欄位
}

func ParseOrders(raw string, allowed map[string]bool, defaultOrder string) []string {
	orders := []string{}
	for _, o := range strings.Split(raw, ",") {
		o = strings.TrimSpace(o)
		if allowed[o] {
			orders = append(orders, o)
		}
	}
	if len(orders) == 0 && defaultOrder != "" {
		orders = append(orders, defaultOrder)
	}
	return orders
}
