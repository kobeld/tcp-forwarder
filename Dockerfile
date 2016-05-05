FROM busybox
ADD main /
CMD ["/main", "-remote", "183.230.96.94:17890"]
EXPOSE 5000
