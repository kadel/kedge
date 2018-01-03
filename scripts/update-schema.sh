#!/usr/bin/env bash

for controller in deployment job deploymentconfig
do 
	schema=$controller"specmod.json"
	p="https://raw.githubusercontent.com/kedgeproject/json-schema/master/master/schema/"$schema
	wget $p
	echo -e 'package validation\n\nvar '${controller^}'specmodJson = `' > $controller.go
	sed -i 's/`//g' $schema
        cat $schema >> $controller.go
	sed -i -e '$a`' $controller.go
	sed -i -e '5s/$/   "additionalProperties": false,/' $controller.go
	mv $controller.go pkg/validation/
	rm $schema

done

