FROM starudream/alpine

WORKDIR /

COPY github-webhook /github-webhook

CMD /github-webhook
