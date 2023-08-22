FROM public.ecr.aws/lambda/go:1

# disables xray when running locally
ENV AWS_XRAY_SDK_DISABLED=TRUE

# run main
CMD [ "main" ]