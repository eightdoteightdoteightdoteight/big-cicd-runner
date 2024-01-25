FROM gobysoft/goby3-debian-build-base
LABEL authors="Brice"

RUN apt install -y git

ENTRYPOINT ["top", "-b"]