FROM public.ecr.aws/lambda/go:1

# disables xray when running locally
ENV AWS_XRAY_SDK_DISABLED=TRUE

# project name arg to build
# build correct path to main.zip
ARG PROJECT_NAME

# run main
CMD [ "main" ]