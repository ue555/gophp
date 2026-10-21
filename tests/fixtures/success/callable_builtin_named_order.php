<?php
function namedBuiltinAdd(int $a, int $b): int { return $a + $b; }
echo call_user_func_array(args: [2, 3], callback: 'namedBuiltinAdd'), '|';
echo call_user_func_array(callback: 'namedBuiltinAdd', args: ['b' => 4, 'a' => 1]), '|';
echo call_user_func(b: 3, callback: 'namedBuiltinAdd', a: 2);
