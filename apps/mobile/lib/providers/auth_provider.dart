import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../core/api.dart';

class AuthState {
  final bool isLoggedIn;
  final String? email;
  final String? tenantId;

  const AuthState({this.isLoggedIn = false, this.email, this.tenantId});

  AuthState copyWith({bool? isLoggedIn, String? email, String? tenantId}) {
    return AuthState(
      isLoggedIn: isLoggedIn ?? this.isLoggedIn,
      email: email ?? this.email,
      tenantId: tenantId ?? this.tenantId,
    );
  }
}

class AuthNotifier extends Notifier<AuthState> {
  @override
  AuthState build() => const AuthState();

  Future<void> login(String email, String password, {String? tenantId}) async {
    await authApi.login(email, password, tenantId: tenantId);
    state = AuthState(
      isLoggedIn: true,
      email: email,
      tenantId: tenantId ?? authApi.tenantId,
    );
  }

  Future<void> register(String email, String password, String fullName) async {
    await authApi.register(email, password, fullName);
  }

  Future<void> logout() async {
    await authApi.logout();
    state = const AuthState();
  }
}

final authProvider = NotifierProvider<AuthNotifier, AuthState>(() => AuthNotifier());
