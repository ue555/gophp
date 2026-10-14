<?php
try {
    strlen([]);
} catch (TypeError $e) {
    echo $e->getMessage(), "\n";
}
