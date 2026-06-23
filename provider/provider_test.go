package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func TestProviderSchema(t *testing.T) {
	ctx := context.Background()
	p := New("test")()

	var resp provider.SchemaResponse
	p.Schema(ctx, provider.SchemaRequest{}, &resp)

	if resp.Schema.Description == "" {
		t.Error("expected provider schema to have a description")
	}

	tokenAttr, ok := resp.Schema.Attributes["token"]
	if !ok {
		t.Fatal("expected provider schema to define a 'token' attribute")
	}

	if !tokenAttr.IsSensitive() {
		t.Error("expected 'token' attribute to be sensitive")
	}
}

func TestProviderMetadata(t *testing.T) {
	ctx := context.Background()
	p := New("1.2.3")()

	var resp provider.MetadataResponse
	p.Metadata(ctx, provider.MetadataRequest{}, &resp)

	if resp.TypeName != "vscale" {
		t.Errorf("expected provider type name 'vscale', got %s", resp.TypeName)
	}

	if resp.Version != "1.2.3" {
		t.Errorf("expected provider version '1.2.3', got %s", resp.Version)
	}
}
