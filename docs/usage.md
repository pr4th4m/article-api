**API usage**

- Create new article

    # Request
    curl -d '@test/post/article1.json' -X POST 'http://localhost:8080/api/v1/articles'
    curl -d '@test/post/article2.json' -X POST 'http://localhost:8080/api/v1/articles'
    curl -d '@test/post/article3.json' -X POST 'http://localhost:8080/api/v1/articles'

    # Response
    {"success":"Article 1 saved successfully"}


- Get article by ID

    # Request
    curl 'http://localhost:8080/api/v1/articles/1'

    # Response
    {
      "id": "1",
      "title": "Artifactory python client",
      "date": "2019-05-10",
      "body": "Artifactory is a artifact repository manager which supports software packages created by different technologies.",
      "tags": [
        "artifactory",
        "python"
      ]
    }


- Search article by tag and date

    # Request
    curl 'http://localhost:8080/api/v1/tags/python/20190510'

    # Response
    {
      "tag": "python",
      "count": 2,
      "articles": [
        "1",
        "2"
      ],
      "related_tags": [
        "artifactory",
        "git"
      ]
    }
