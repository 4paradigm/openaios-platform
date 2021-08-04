const path = require('path');

module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'node',
  globals: {
    'ts-jest': {
        isolatedModules: true
    }
  },
  verbose: true,
  testURL: 'http://localhost/',
  rootDir: path.resolve(__dirname, './'),
  moduleFileExtensions: [
      'js',
      'json',
      'tsx',
      'ts'
  ],
  moduleNameMapper: {
      '^@\/(.*?\.?(ts|js|tsx)?|)$': '<rootDir>/src/$1',
      '\\.(jpg|jpeg|png|gif|eot|otf|webp|svg|ttf|woff|woff2|mp4|webm|wav|mp3|m4a|aac|oga)$': '<rootDir>/tests/mocks/fileMock.js', // 模拟加载静态文件
      '\\.(css|less|scss|sass)$': '<rootDir>/tests/mocks/styleMock.js'　　// 模拟加载样式文件   
  },
  transform: {
      '^.+\\.js$': '<rootDir>/node_modules/ts-jest',
      '^.+\\.ts$': '<rootDir>/node_modules/ts-jest',
      '.*\\.(tsx)$': '<rootDir>/node_modules/ts-jest',
      '^.+\\.(jpg|jpeg|png|gif|eot|otf|webp|svg|ttf|woff|woff2|mp4|webm|wav|mp3|m4a|aac|oga)$':
      'jest-transform-stub',
  },
  testMatch: [
      '<rootDir>/src/**/__tests__/*.test.js',
      '<rootDir>/src/**/__tests__/*.test.tsx'
  ],
  transformIgnorePatterns: ["node_modules/(?!(cess-ui|smart-chart)/)"],   
  setupFiles: ['<rootDir>/tests/setup.ts'],
  setupFilesAfterEnv: ['<rootDir>/tests/setup.ts'],
  snapshotSerializers: ['enzyme-to-json/serializer'],
  coverageDirectory: '<rootDir>/tests/unit/coverage', // 覆盖率报告的目录
  collectCoverageFrom: [
      'src/components/**/*.{ts,tsx}',
      '!src/demos/**/*.{ts,tsx}',
      '!src/index.ts',
      '!lib',
      '!node_modules'
  ]
};