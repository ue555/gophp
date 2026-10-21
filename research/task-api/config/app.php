<?php

return [
    'name' => 'GoPHP Task API',
    'env' => env('APP_ENV', 'testing'),
    'debug' => false,
    'url' => env('APP_URL', 'http://127.0.0.1:8091'),
    'timezone' => 'UTC',
    'locale' => 'en',
    'fallback_locale' => 'en',
    'key' => env('APP_KEY', 'base64:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA='),
    'cipher' => 'AES-256-CBC',
    'maintenance' => ['driver' => 'file'],
];
