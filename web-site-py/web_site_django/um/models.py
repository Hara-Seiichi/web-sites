from django.db import models
from sqlalchemy import false, true

class SexModel(models.Model):
    sex = models.CharField(max_length=1,default="")
    def __str__(self):
        return self.sex

class UsersModel(models.Model):
    userid = models.CharField(max_length=13, default="")
    name = models.CharField(max_length=50, default="")
    age = models.PositiveIntegerField(default=0)

    #nullを許容しないと既存データがある場合マイグレーションが失敗する。
    sex = models.ForeignKey(SexModel, on_delete=models.PROTECT, null=True)
