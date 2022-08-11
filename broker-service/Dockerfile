FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/brokerService /app

CMD [ "/app/brokerService" ]