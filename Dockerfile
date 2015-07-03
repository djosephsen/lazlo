FROM scratch
ADD dockerdep/ca-certificates.crt /etc/ssl/certs/
ADD lazlo /
CMD ["/lazlo"]
