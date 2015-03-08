FROM scratch
ADD docker-demo /bin/docker-demo
ADD static /static
ADD templates /templates
EXPOSE 8080
ENTRYPOINT ["/bin/docker-demo"]
CMD ["-listen=:8080"]
