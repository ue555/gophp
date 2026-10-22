<?php
namespace Probe {
    use function getmypid as process_id;
    $pid = getmypid();
    echo gettype($pid), '|', ($pid > 0 ? 'positive' : 'invalid'), "\n";
    echo ($pid === \getmypid() && $pid === process_id() && $pid === GETMYPID()) ? "stable\n" : "different\n";
    $call = 'getmypid';
    echo ($call() === $pid && call_user_func('getmypid') === $pid) ? "dynamic\n" : "different\n";
    echo function_exists('getmypid') ? "exists\n" : "missing\n";
    echo function_exists(__NAMESPACE__.'\\getmypid') ? "wrong namespace\n" : "fallback\n";
    try { namespace\getmypid(); } catch (\Error $e) { echo $e->getMessage(), "\n"; }
}
namespace Shadow {
    function getmypid() { return 'local'; }
    echo getmypid(), '|', (\getmypid() > 0 ? 'global' : 'invalid'), "\n";
}
