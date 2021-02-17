FROM scratch

COPY tcli .

ENTRYPOINT [ "/tcli" ]
CMD ["--help"]

