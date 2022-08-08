from django.contrib import messages
from django.shortcuts import render, redirect
from django.contrib.auth.models import User
from django.db import IntegrityError
from django.db.models import Q
from django.contrib.auth import authenticate, login, logout
from django.urls import reverse_lazy

from django.views.generic import CreateView, DeleteView, DetailView, ListView, UpdateView
from um.forms import UserForm

from um.models import UsersModel

def Signup(request):
    if(request.method == 'POST'):
        username = request.POST['username']
        password = request.POST['password']
        try:
            # 自前で作成したユーザ管理テーブルではなくadmin/の管理下
            user = User.objects.create_user(username, '', password)
            return redirect('signin')
        except IntegrityError:
            return render(request, 'signup.html', {'error': 'The user is already registered.'})
        except ValueError:
            # DBアクセスの前に値検証はする方が良い。今回は趣味範囲としておｋにする。。
            return render(request, 'signup.html', {'error': 'Please enter both your username and password.'})
    return render(request, 'signup.html', {})

def Signin(request):
    if request.method == 'POST':
        username = request.POST['username']
        password = request.POST['password']
        # 作成したアプリと紛らわしいが、Djangoに標準機能でログイン認下証は行う。
        # ※自前で作成したユーザ管理テーブルではなくadmin/の管理下
        user = authenticate(username=username, password=password)
        if user is not None:
            login(request, user)
            return redirect('list')
        else:
            # レスポンスされた画面では「error」の名称で参照出来る
            return render(request, 'signin.html', {'error': 'That user does not exist.'})
    return render(request, 'signin.html', {})

def Signout(request):
    # リクエストの認証情報、セッション情報をお掃除してくれる。
    logout(request)
    return redirect('signin')

class UserListView(ListView):
    # django.views.generic.ListViewを継承する時は
    # 「model」にテーブルのデータを入れるのはお約束
    model = UsersModel
    template_name = 'list.html'

    def get_context_data(self):
        # 自画面遷移で再表示したときにformの入力値が無くならない様にする
        # 登録、更新、削除、詳細から戻ってきた場合、今回は検索条件の保持はしていない。
        ctx = super().get_context_data()
        ctx['userid'] = '' if self.request.GET.get('userid') is None else self.request.GET.get('userid')
        ctx['name'] = '' if self.request.GET.get('name') is None else self.request.GET.get('name')
        ctx['age'] = '' if self.request.GET.get('age') is None else self.request.GET.get('age')
        ctx['sex'] = '' if self.request.GET.get('sex') is None else self.request.GET.get('sex')
        return ctx

    def get_queryset(self, **kwargs):
        queryset = super().get_queryset(**kwargs)
        user_id = self.request.GET.get('userid')
        user_name = self.request.GET.get('name')
        user_age = self.request.GET.get('age')
        user_sex = self.request.GET.get('sex')

        # 検索条件
        # 空のQオブジェクトを作成しておく
        condition_id = Q()
        condition_name = Q()
        condition_age = Q()
        condition_sex = Q()

        # 入力があれば条件に加える
        if user_id:
            condition_id = Q(userid__icontains=user_id)
        if user_name:
            condition_name = Q(name__icontains=user_name)
        if user_age:
            condition_age = Q(age__exact=user_age)
        if user_sex:
            condition_sex = Q(sex__sex__exact=user_sex)

        # Qオブジェクトに値が設定されたものだけFilterが効く
        queryset = UsersModel.objects.filter(
                                        condition_id &
                                        condition_name &
                                        condition_age &
                                        condition_sex
                                        ).order_by("pk")

        return queryset

class UserDetailView(DetailView):
    # ２行でquerystringに渡したPKで検索出来る。凄い。
    model = UsersModel 
    template_name = 'detail.html'

class UserCreateView(CreateView):
    model = UsersModel 
    template_name = 'create.html'
    form_class = UserForm

    # 成功した時のURL(定義しないと落ちる)
    success_url = reverse_lazy('search')

    # バリデーションでエラーになったとき
    # 今回はまとめてエラーをリスト形式で返却。
    # 画面に戻った時のレイアウト調整が大変だったのでモーダル表示で逃げ。。
    def form_invalid(self, form):
        messages.add_message(self.request, messages.ERROR, form.errors)
        return redirect('create')

class UserUpdateView(UpdateView):
    # 継承元を変えるだけで登録の処理と変わらない。凄い。
    # querystringにPKを渡すのを忘れずに。
    model = UsersModel 
    template_name = 'update.html'
    form_class = UserForm

    # 成功した時のURL(定義しないと落ちる)
    success_url = reverse_lazy('search')

    # バリデーションでエラーになったとき
    def form_invalid(self, form):
        messages.add_message(self.request, messages.ERROR, form.errors)
        return redirect('update')

class UserDeleteView(DeleteView):
    # querystringにPKを渡すのを忘れずに。これだけ凄いｗ
    model = UsersModel 
    template_name = 'delete.html'

        # 成功した時のURL
    success_url = reverse_lazy('search')

