# terraform-deployer

Repository and Github Action for updating and releasing in Terraform Cloud

## Usage

If you have a Terraform Workflow with the variable `my_docker_image`, and you would like to update the value on each commit to your repository, you can create a `.github/workflows/release.yml` that looks something like this:

```yaml
name: Deploy
on:
  push:
    branches:
      - main
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Build and push
        id: push
        run: |
          docker build -t my-repo/my-app:latest .
          DOCKER_IMAGE="my-repo/my-app@$(docker push -q my-repo/my-app:latest)"
          echo "docker_image_value=${DOCKER_IMAGE}" >> $GITHUB_OUTPUT

      - name: Deploy to dev
        uses: xmtplabs/terraform-deployer@v1
        with:
          terraform-token: ${{ secrets.TERRAFORM_TOKEN }}
          terraform-org: ${{ secrets.TERRAFORM_ORG }}
          terraform-workspace: dev
          variable-name: my_docker_image
          variable-value: ${{ steps.push.outputs.docker_image_value }}
          variable-value-required-prefix: my-repo/my-app@
```
