package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateSSHKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected POST request, got %s", req.Method)
		}
		if req.URL.Path != "/sshkeys" {
			t.Errorf("expected path /sshkeys, got %s", req.URL.Path)
		}
		if req.Header.Get("X-Token") != "test-token" {
			t.Errorf("expected X-Token header test-token, got %s", req.Header.Get("X-Token"))
		}

		response := []SSHKey{
			{ID: 123, Name: "test-key", Key: "ssh-rsa AAA..."},
		}
		respBytes, _ := json.Marshal(response)
		rw.WriteHeader(http.StatusOK)
		rw.Write(respBytes)
	}))
	defer server.Close()

	c := NewClient("test-token")
	c.BaseURL = server.URL

	key, err := c.CreateSSHKey("test-key", "ssh-rsa AAA...")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if key.ID != 123 {
		t.Errorf("expected ID 123, got %d", key.ID)
	}
	if key.Name != "test-key" {
		t.Errorf("expected name test-key, got %s", key.Name)
	}
}

func TestGetScalet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			t.Errorf("expected GET request, got %s", req.Method)
		}
		if req.URL.Path != "/scalets/10299" {
			t.Errorf("expected path /scalets/10299, got %s", req.URL.Path)
		}

		// GET /scalets/:ctid returns an array with one element in VScale API
		response := []Scalet{
			{
				CTID:     10299,
				Name:     "Hollow-Star",
				Status:   "started",
				Location: "spb0",
				RPlan:    "large",
			},
		}
		respBytes, _ := json.Marshal(response)
		rw.WriteHeader(http.StatusOK)
		rw.Write(respBytes)
	}))
	defer server.Close()

	c := NewClient("test-token")
	c.BaseURL = server.URL

	scalet, err := c.GetScalet(10299)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if scalet.CTID != 10299 {
		t.Errorf("expected CTID 10299, got %d", scalet.CTID)
	}
	if scalet.Name != "Hollow-Star" {
		t.Errorf("expected name Hollow-Star, got %s", scalet.Name)
	}
	if scalet.Status != "started" {
		t.Errorf("expected status started, got %s", scalet.Status)
	}
}

func TestCreateDomain(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			t.Errorf("expected POST request, got %s", req.Method)
		}
		if req.URL.Path != "/domains" {
			t.Errorf("expected path /domains, got %s", req.URL.Path)
		}

		response := Domain{
			ID:   36,
			Name: "example.com",
		}
		respBytes, _ := json.Marshal(response)
		rw.WriteHeader(http.StatusOK)
		rw.Write(respBytes)
	}))
	defer server.Close()

	c := NewClient("test-token")
	c.BaseURL = server.URL

	domain, err := c.CreateDomain("example.com", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if domain.ID != 36 {
		t.Errorf("expected ID 36, got %d", domain.ID)
	}
	if domain.Name != "example.com" {
		t.Errorf("expected name example.com, got %s", domain.Name)
	}
}
