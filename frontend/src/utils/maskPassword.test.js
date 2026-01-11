/**
 * 密码脱敏工具函数测试
 */

import { maskKnownPasswords } from './maskPassword';

describe('maskPassword', () => {
  describe('maskKnownPasswords', () => {
    it('应该替换精确匹配的密码', () => {
      const text = '1. 打开登录页面\n2. 输入用户名"root"\n3. 输入密码"root123"\n4. 点击登录按钮';
      const knownPasswords = ['root123'];
      const result = maskKnownPasswords(text, knownPasswords);
      
      expect(result).toContain('******');
      expect(result).not.toContain('root123');
      expect(result).toContain('root'); // 用户名不应被替换
    });

    it('应该替换多个不同的密码', () => {
      const text = 'password1: admin123, password2: test456';
      const knownPasswords = ['admin123', 'test456'];
      const result = maskKnownPasswords(text, knownPasswords);
      
      expect(result).not.toContain('admin123');
      expect(result).not.toContain('test456');
      expect(result.match(/\*\*\*\*\*\*/g)).toHaveLength(2);
    });

    it('应该处理JSON格式的密码', () => {
      const text = '{"username": "root", "password": "root123"}';
      const knownPasswords = ['root123'];
      const result = maskKnownPasswords(text, knownPasswords);
      
      expect(result).toContain('******');
      expect(result).not.toContain('root123');
    });

    it('应该处理Playwright脚本中的密码', () => {
      const text = `await page.getByPlaceholder('密码').fill('root123');`;
      const knownPasswords = ['root123'];
      const result = maskKnownPasswords(text, knownPasswords);
      
      expect(result).toContain('******');
      expect(result).not.toContain('root123');
    });

    it('应该按密码长度降序处理，避免短密码影响长密码', () => {
      const text = 'password: root123456';
      const knownPasswords = ['root123', 'root123456']; // 短的在前
      const result = maskKnownPasswords(text, knownPasswords);
      
      // 应该替换长密码，而不是只替换短密码部分
      expect(result).toBe('password: ******');
    });

    it('空文本应返回原值', () => {
      expect(maskKnownPasswords('', ['pass'])).toBe('');
      expect(maskKnownPasswords(null, ['pass'])).toBe(null);
      expect(maskKnownPasswords(undefined, ['pass'])).toBe(undefined);
    });

    it('空密码列表应返回原文本', () => {
      const text = 'password: root123';
      expect(maskKnownPasswords(text, [])).toBe(text);
      expect(maskKnownPasswords(text, null)).toBe(text);
      expect(maskKnownPasswords(text, undefined)).toBe(text);
    });

    it('应该转义正则特殊字符', () => {
      const text = 'password: root$123.test';
      const knownPasswords = ['root$123.test'];
      const result = maskKnownPasswords(text, knownPasswords);
      
      expect(result).toBe('password: ******');
    });

    it('应该替换所有出现的密码', () => {
      const text = '步骤1: 输入root123\n步骤2: 确认root123\n步骤3: 验证root123';
      const knownPasswords = ['root123'];
      const result = maskKnownPasswords(text, knownPasswords);
      
      expect(result.match(/\*\*\*\*\*\*/g)).toHaveLength(3);
      expect(result).not.toContain('root123');
    });
  });
});
