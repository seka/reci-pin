export const VALIDATION_RULES = {
    EMAIL: {
        MAX_LENGTH: 200,
    },
    PASSWORD: {
        MIN_LENGTH: 8,
        MAX_LENGTH: 200,
    },
    RECIPE: {
        NAME_MAX_LENGTH: 100,
        MEMO_MAX_LENGTH: 1000,
    },
    TAG: {
        NAME_MAX_LENGTH: 30,
    },
} as const;
