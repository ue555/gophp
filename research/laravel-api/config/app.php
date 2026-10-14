<?php

return [
    'name' => 'GoPHP Laravel API probe',
    'env' => env('APP_ENV', 'testing'),
    'debug' => false,
    'url' => env('APP_URL', 'http://127.0.0.1:8080'),
    'timezone' => 'UTC',
    'locale' => 'en',
    'fallback_locale' => 'en',
    'key' => env('APP_KEY', 'base64:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA='),
    'cipher' => 'AES-256-CBC',
    'maintenance' => ['driver' => 'file'],
];
