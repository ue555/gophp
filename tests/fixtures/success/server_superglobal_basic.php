<?php
function serverMissing(): bool {
    return isset($_SERVER['gophp_missing_2026']);
}
echo gettype($_SERVER), '|', serverMissing() ? 'wrong' : 'clear';
