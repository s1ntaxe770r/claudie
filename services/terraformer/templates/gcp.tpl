{{- $clusterName := .ClusterName}}
{{- $clusterHash := .ClusterHash}}
{{$index :=  0}}

provider "google" {
  credentials = "${file("{{(index .NodePools $index).Provider.SpecName}}")}"
  region      = "{{(index .NodePools 0).Region}}"
  project     = "{{(index .NodePools 0).Provider.GcpProject}}"
  alias       = "k8s-nodepool"
}

resource "google_compute_network" "network" {
  provider                = google.k8s-nodepool
  name                    = "{{ $clusterName }}-{{ $clusterHash }}-network"
  auto_create_subnetworks = false
}

resource "google_compute_firewall" "firewall" {
  provider     = google.k8s-nodepool
  name         = "{{ $clusterName }}-{{ $clusterHash }}-firewall"
  network      = google_compute_network.network.self_link

  allow {
    protocol = "UDP"
    ports    = ["51820"]
  }

  {{ if index .Metadata "loadBalancers" | targetPorts | isMissing 6443 }}
  allow {
      protocol = "TCP"
      ports    = ["6443"]
  }
  {{ end }}

  allow {
      protocol = "TCP"
      ports    = ["22"]
  }

  allow {
      protocol = "icmp"
   }

  source_ranges = [
      "0.0.0.0/0",
   ]
}

{{range $i, $nodepool := .NodePools}}
resource "google_compute_subnetwork" "{{$nodepool.Name}}-subnet" {
  provider      = google.k8s-nodepool
  name          = "{{ $nodepool.Name }}-{{ $clusterHash }}-subnet"
  network       = google_compute_network.network.self_link
  region        = "{{$nodepool.Region}}"
  ip_cidr_range = "{{getCIDR "10.0.0.0/24" 2 $i}}"
}

resource "google_compute_instance" "{{ $nodepool.Name }}" {
  provider     = google.k8s-nodepool
  count        = {{ $nodepool.Count }}
  zone         = "{{$nodepool.Zone}}"
  name         = "{{ $clusterName }}-{{ $clusterHash }}-{{ $nodepool.Name }}-${count.index + 1}"
  machine_type = "{{ $nodepool.ServerType }}"
  allow_stopping_for_update = true
  boot_disk {
    initialize_params {
      size = "{{ $nodepool.DiskSize }}"
      image = "{{ $nodepool.Image }}"
    }
  }
  network_interface {
    subnetwork = google_compute_subnetwork.{{$nodepool.Name}}-subnet.self_link
    access_config {}
  }
  metadata = {
    ssh-keys = "root:${file("./public.pem")}"
  }
  metadata_startup_script = "echo 'PermitRootLogin without-password' >> /etc/ssh/sshd_config && echo 'PubkeyAuthentication yes' >> /etc/ssh/sshd_config && service sshd restart"
}

output "{{ $nodepool.Name }}" {
  value = {
    for node in google_compute_instance.{{ $nodepool.Name }}:
    node.name => node.network_interface.0.access_config.0.nat_ip
  }
}
{{end}}


