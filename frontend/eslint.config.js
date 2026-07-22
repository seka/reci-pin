// @ts-check
const eslint = require("@eslint/js");
const tseslint = require("typescript-eslint");
const angular = require("angular-eslint");
const eslintConfigPrettier = require("eslint-config-prettier");
// For more info, see https://github.com/storybookjs/eslint-plugin-storybook#configuration-flat-config-format
const storybook = require("eslint-plugin-storybook");

module.exports = tseslint.config(
    {
        files: ["**/*.ts"],
        extends: [
            eslint.configs.recommended,
            ...tseslint.configs.recommended,
            ...tseslint.configs.stylistic,
            ...angular.configs.tsRecommended,
            eslintConfigPrettier,
        ],
        processor: angular.processInlineTemplates,
        rules: {
            "@typescript-eslint/no-unused-vars": ["error", { argsIgnorePattern: "^_", varsIgnorePattern: "^_" }],
            "@typescript-eslint/no-inferrable-types": "off",
            // Angular 22 で OnPush がデフォルト化され、ng update が既存挙動維持のため
            // 全コンポーネントに ChangeDetectionStrategy.Eager を付与した。OnPush への
            // 実移行は別タスクとして扱うため、このルールは一時的に無効化する。
            "@angular-eslint/prefer-on-push-component-change-detection": "off",
        },
    },
    {
        files: ["**/*.html"],
        extends: [
            ...angular.configs.templateRecommended,
            ...angular.configs.templateAccessibility,
            eslintConfigPrettier,
        ],
        rules: {},
    },
    ...storybook.configs["flat/recommended"]
);
