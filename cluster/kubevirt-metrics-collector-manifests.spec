#debuginfo not supported with Go
%global debug_package %{nil}
%global source_name v%{version}

%global commit 6e4ee92d5cc434186c9843e3c014aff346a18714
%global build_hash %(c=%{commit}; echo ${c:0:7})
%global spec_release 1

Name:		kubevirt-metrics-collector-manifests
Version:        0.12.0.5
Release:	%{spec_release}.%{build_hash}
Summary:	manifests to deploy kubevirt-metrics-collector
		
License:	ASL 2.0
URL:		https://github.com/fromanirh/kubevirt-metrics-collector
Source0:	config-map.yaml
Source1:	k8s-daemonset.yaml
Source2:	okd-account-scc.yaml
Source3:	okd-daemonset.yaml
Source4:	vmi-service-monitor.yaml
Source5:	vmi-service.yaml
Source6:	fake-node-exporter/config-map.yaml
Source7:	fake-node-exporter/okd-daemonset.yaml

BuildArch: noarch
%description
manifests to deploy the kubevirt metrics collector

%prep

%build

%install
rm -rf  %{buildroot}
install -d -m 0755 %{buildroot}%{_datadir}/%{name}/manifests
install -d -m 0755 %{buildroot}%{_datadir}/%{name}/manifests/fake-node-exporter
cp -v %{SOURCE0} %{buildroot}%{_datadir}/%{name}/manifests/config-map.yaml
cp -v %{SOURCE1} %{buildroot}%{_datadir}/%{name}/manifests/k8s-daemonset.yaml
cp -v %{SOURCE2} %{buildroot}%{_datadir}/%{name}/manifests/okd-account-scc.yaml
cp -v %{SOURCE3} %{buildroot}%{_datadir}/%{name}/manifests/okd-daemonset.yaml
cp -v %{SOURCE4} %{buildroot}%{_datadir}/%{name}/manifests/vmi-service-monitor.yaml
cp -v %{SOURCE5} %{buildroot}%{_datadir}/%{name}/manifests/vmi-service.yaml
cp -v %{SOURCE6} %{buildroot}%{_datadir}/%{name}/manifests/fake-node-exporter/config-map.yaml
cp -v %{SOURCE7} %{buildroot}%{_datadir}/%{name}/manifests/fake-node-exporter/okd-daemonset.yaml

%files 
%{_datadir}/%{name}/manifests/


%changelog
* Wed Jan 9 2019 Francesco Romani <fromani@redhat.com> - 0.12.0.5-6e4ee92d
- 0.12.0 Release with compatibility with OKD.

* Thu Dec 17 2018 Francesco Romani <fromani@redhat.com> - 0.12.0.2-d8a6a6a5
- 0.12.0 Release requiring kubevirt >= 0.12.0.

* Thu Dec 17 2018 Francesco Romani <fromani@redhat.com> - 0.12.0.1-af1593f1
- 0.12.0 Release

