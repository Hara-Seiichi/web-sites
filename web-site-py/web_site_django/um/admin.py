from django.contrib import admin

from um.models import SexModel, UsersModel

# models.pyで作成したテーブルはここに記載する。
# Djangoの管理者用画面(admin)からデータの追加、変更、削除などが出来る
admin.site.register(SexModel)
admin.site.register(UsersModel)
