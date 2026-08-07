import 'package:dio/dio.dart';

const _base = String.fromEnvironment(
  'API_URL',
  defaultValue: 'http://10.0.2.2:3001', // Android emulator → localhost
);

class _TokenStore {
  String? accessToken;
  String? refreshToken;
  String? tenantId;
}

final _tokens = _TokenStore();

final _dio = Dio(BaseOptions(baseUrl: _base))
  ..interceptors.add(InterceptorsWrapper(
    onRequest: (options, handler) {
      if (_tokens.accessToken != null) {
        options.headers['Authorization'] = 'Bearer ${_tokens.accessToken}';
      }
      handler.next(options);
    },
  ));

class AuthApi {
  Future<void> login(String email, String password, {String? tenantId}) async {
    final res = await _dio.post('/auth/login', data: {
      'email': email,
      'password': password,
      if (tenantId case final id?) 'tenant_id': id,
    });
    _tokens.accessToken = res.data['access_token'];
    _tokens.refreshToken = res.data['refresh_token'];
    _tokens.tenantId = tenantId;
  }

  Future<void> register(String email, String password, String fullName) async {
    await _dio.post('/auth/register', data: {
      'email': email,
      'password': password,
      'full_name': fullName,
    });
  }

  void logout() {
    _tokens.accessToken = null;
    _tokens.refreshToken = null;
    _tokens.tenantId = null;
  }

  String? get accessToken => _tokens.accessToken;
  String? get tenantId => _tokens.tenantId;
}

class TenantsApi {
  Future<List<Map<String, dynamic>>> list() async {
    final res = await _dio.get('/tenants');
    return List<Map<String, dynamic>>.from(res.data);
  }

  Future<Map<String, dynamic>> create(String name, String slug) async {
    final res = await _dio.post('/tenants', data: {'name': name, 'slug': slug});
    return Map<String, dynamic>.from(res.data);
  }
}

class PowerSyncApi {
  Future<String> getToken() async {
    final res = await _dio.get('/powersync/token');
    return res.data['token'] as String;
  }
}

final authApi = AuthApi();
final tenantsApi = TenantsApi();
final powerSyncApi = PowerSyncApi();
