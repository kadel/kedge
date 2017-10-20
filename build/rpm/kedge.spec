Summary:       Simple, Concise & Declarative Kubernetes Applications http://kedgeproject.org
Name:          kedge
Version:       0.3.0
Release:       1%{?dist}
License:       Apache License 2.0

%global github_project kedgeproject
%global import_path github.com/%{github_project}/%{name}


Url:           https://github.com/%{github_project}/%{name}
Source0:       https://github.com/%{github_project}/%{name}/archive/v%{version}.tar.gz

BuildRequires: golang
BuildRequires: make

%global _dwz_low_mem_die_limit 0

%define gobuild(o:) go build -ldflags "${LDFLAGS:-} -B 0x$(head -c20 /dev/urandom|od -An -tx1|tr -d ' \\n')" %{?**};


%description
Kedge is a simple, easy and declarative way to define and deploy applications to Kubernetes by writing very concise application definitions.

%prep
%setup -q -n %{name}-%{version}

%build
mkdir -p ./_build/src/github.com/%{github_project}/
ln -s $(pwd) ./_build/src/%{import_path}
export GOPATH=$(pwd)/_build/
%gobuild %{import_path}

%install
install -d %{buildroot}%{_bindir}
install -p -m 0755 ./kedge %{buildroot}%{_bindir}/kedge

%files
%{_bindir}/kedge

%clean
rm -rf %{buildroot}


%changelog
* Mon Oct 23 2017 Tomas Kral <tkral@redhat.com> 0.3.0-1
- initial version
