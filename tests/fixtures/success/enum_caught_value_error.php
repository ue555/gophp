<?php
enum Phase: int { case Start = 1; }
try {
    Phase::from(9);
    echo 'bad';
} catch (ValueError $error) {
    echo 'caught|';
} finally {
    echo 'finally';
}
