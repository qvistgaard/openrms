#!/usr/bin/env bash
export BUNDLE_GEMFILE=docs/site/Gemfile
bundle install
bundle exec jekyll build --config docs/site/config.yaml --source docs/site --destination _site
bundle exec jekyll serve --config docs/site/config.yaml --source docs/site --destination _site --host=0.0.0.0