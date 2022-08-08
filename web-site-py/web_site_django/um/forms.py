from django import forms

from .models import UsersModel

class UserForm(forms.ModelForm):
    userid = forms.CharField(required=True)
    name = forms.CharField(required=True)
    age = forms.IntegerField(required=True, min_value=0)
    sex = forms.IntegerField(required=True)
    class Meta:
        # どのモデルをフォームにするか指定
        model = UsersModel
        # そのフォームの中から表示するフィールドを指定
        fields = ('userid', 'name', 'age', 'sex')

        # エラーメッセージをまとめて書ける
        error_messages = {
            'userid': {
                'required': 'input required!',
            },
            'name': {
                'required': 'input required!',
            },
            'age': {
                'required': 'input required!',
                'min_value': 'There are no ages less than zero.',
            },
            'sex': {
                'required': 'input required!',
            },
        }