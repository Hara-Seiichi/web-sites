#
# 秘匿すべき設定情報を記載する。
# 

# プロジェクト生成時に自動生成される。
SECRET_KEY = ''

# postgresqlを利用する場合の設定
# psycopg2を実行する仮想環境へインストールする必要がある。
DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.postgresql',
        'NAME': 'postgres',
        'USER': 'user',
        'PASSWORD': 'pass',
        'HOST': '127.0.0.1',
        'PORT': '5432',
    }
}