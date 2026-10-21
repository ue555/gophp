<?php declare(strict_types=1);
$mapped = array_map(static fn (string $key): string => $key, [70616765]);
$filtered = array_filter([5], static fn (string $value): bool => $value === '5');
echo $mapped[0], '|', count($filtered), '|';
echo call_user_func(static fn (string $value): string => $value, 7), '|';
echo call_user_func_array(static fn (string $value): string => $value, [8]), '|';
try {
    (static fn (string $value): string => $value)(9);
} catch (TypeError $error) {
    echo 'strict';
}
