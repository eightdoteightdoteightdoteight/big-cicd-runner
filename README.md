# big-cicd-runner
Runner for the big-cicd project, the runner deploys the app on our Kubernetes environment. \
If triggered from our app it just pulls an already existing version of our application (only CD). \
If triggered from Github it builds a new version of the app and performs the commands defined in the big_ci.yml file (both CI and CD).

CI includes the build of the application and the other defined commands of the big_ci.yml file. \
CD includes deployment on Kuebrnetes, pentest and healthcheck. \

This project is entirely built in Go.
