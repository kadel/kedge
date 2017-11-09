Summary:       Simple, Concise & Declarative Kubernetes Applications http://kedgeproject.org
Name:          kedge
Version:       0.4.0
Release:       1%{?dist}
License:       Apache License 2.0

%global github_project kedgeproject
%global import_path github.com/%{github_project}/%{name}


Url:           https://github.com/%{github_project}/%{name}
Source0:       https://github.com/%{github_project}/%{name}/archive/v%{version}.tar.gz

BuildRequires: golang
BuildRequires: make

%global _dwz_low_mem_die_limit 0

%description
Kedge is a simple, easy and declarative way to define and deploy applications to Kubernetes by writing very concise application definitions.

%prep
%setup -q -n %{name}-%{version}

%build
mkdir -p ./_build/src/github.com/%{github_project}/
ln -s $(pwd) ./_build/src/%{import_path}
export GOPATH=$(pwd)/_build/
make bin

%install
install -d %{buildroot}%{_bindir}
install -p -m 0755 ./kedge %{buildroot}%{_bindir}/kedge

%files
%{_bindir}/kedge

%clean
rm -rf %{buildroot}


%changelog
* Thu Nov 09 2017 Tomas Kral <tkral@redhat.com> 0.4.0-1
- update to v0.4.0

* Mon Oct 23 2017 Tomas Kral <tkral@redhat.com> 0.3.0-1
- initial version
