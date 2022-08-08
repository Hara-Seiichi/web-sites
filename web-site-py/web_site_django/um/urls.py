from um import views
from django.urls import path

from um.views import UserCreateView, UserDeleteView, UserListView, UserDetailView, UserUpdateView

urlpatterns = [
    path('', views.Signup, name='signup'),
    path('signin/', views.Signin, name='signin'),
    path('signout/', views.Signout, name='signout'),
    path('list/', UserListView.as_view(), name='list'),
    path('search/', UserListView.as_view(), name='search'),
    path('detali/<int:pk>/', UserDetailView.as_view(), name='detali'),
    path('create/', UserCreateView.as_view(), name='create'),
    path('update/<int:pk>/', UserUpdateView.as_view(), name='update'),
    path('delete/<int:pk>/', UserDeleteView.as_view(), name='delete'),
]
