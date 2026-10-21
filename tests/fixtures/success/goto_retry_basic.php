<?php
function retryOnce() {
    $attempts = 0;
    beginning:
    $attempts++;
    try {
        if ($attempts === 1) throw new Exception('retry');
        return $attempts;
    } catch (Exception $e) {
        goto beginning;
    }
}
echo retryOnce();
