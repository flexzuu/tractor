# This Dockerfile is not for use so much as to verify a working dev environment
# from scratch. Building this with `make test-env` will show missing dependencies
# and help clarify what software is necessary for a new development setup.

FROM golang:1.13.6-alpine3.10
RUN echo "@community http://nl.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
RUN apk --update add nodejs yarn curl make pkgconf gtk+3.0-dev webkit2gtk-dev libappindicator-dev@community g++ python bash git ca-certificates
RUN yarn global add typescript
RUN git clone https://github.com/manifold/tractor
RUN make -C tractor versions setup