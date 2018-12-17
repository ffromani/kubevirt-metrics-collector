#debuginfo not supported with Go
%global debug_package %{nil}
%global source_name v%{version}

%global commit 6b02f1500062ec8fb625db22584820bb8edd8afb
%global build_hash %(c=%{commit}; echo ${c:0:7})
%global spec_release 1

Name:		kubevirt-metrics-collector-manifests
Version:        0.10.0
Release:	%{spec_release}.%{build_hash}
Summary:	manifests to deploy kubevirt-metrics-collector
		
License:	ASL 2.0
URL:		https://github.com/fromanirh/kubevirt-metrics-collector
Source0:	kubevirt-metrics-collector-config-map.yaml
Source1:	kubevirt-metrics-collector-k8s.yaml
Source2:	kubevirt-metrics-collector-ocp.yaml
Source3:	kubevirt-service-monitor-vmi.yaml
Source4:	kubevirt-service-vmi.yaml

BuildArch: noarch
%description
manifests to deploy the kubevirt metrics collector

%prep

%build

%install
rm -rf  %{buildroot}
install -d -m 0755 %{buildroot}%{_datadir}/%{name}/manifests
cp -v %{SOURCE0} %{buildroot}%{_datadir}/%{name}/manifests/kubevirt-metrics-collector-config-map.yaml
cp -v %{SOURCE1} %{buildroot}%{_datadir}/%{name}/manifests/kubevirt-metrics-collector-k8s.yaml
cp -v %{SOURCE2} %{buildroot}%{_datadir}/%{name}/manifests/kubevirt-metrics-collector-ocp.yaml
cp -v %{SOURCE3} %{buildroot}%{_datadir}/%{name}/manifests/kubevirt-service-monitor-vmi.yaml
cp -v %{SOURCE4} %{buildroot}%{_datadir}/%{name}/manifests/kubevirt-service-vmi.yaml

%files 
%{_datadir}/%{name}/manifests/


%changelog
* Thu Dec 17 2018 Francesco Romani <fromani@redhat.com> - 0.10.0.4-6b02f150
- 0.10.0 Release

