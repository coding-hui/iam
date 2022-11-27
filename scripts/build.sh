#!/bin/bash

#./mvnw clean install -B ${@}
./mvnw clean install -B -DskipTests ${@}