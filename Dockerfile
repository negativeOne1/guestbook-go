FROM scratch
WORKDIR /app
COPY /bin/guestbook .
COPY public .
CMD ["/app/guestbook"]
EXPOSE 8000
