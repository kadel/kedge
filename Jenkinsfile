#!/usr/bin/groovy
@Library('github.com/fabric8io/fabric8-pipeline-library@master')
def dummy
goNode{
  dockerNode{
    ws{
      if (env.BRANCH_NAME.startsWith('PR-')) {
        goCI{
          githubOrganisation = 'kedgeproject'
          dockerOrganisation = 'kedgeproject'
          project = 'kedge'
          makeTarget = 'clean test cross'
        }
      } 
    }
  }
}
