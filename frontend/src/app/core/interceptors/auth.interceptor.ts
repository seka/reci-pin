import { HttpInterceptorFn } from '@angular/common/http';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  // ブラウザが Cookie (auth_token) を自動的に送信するため、
  // 手動での Authorization ヘッダー付与は不要になりました。
  return next(req);
};
