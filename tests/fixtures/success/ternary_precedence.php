<?php
echo (false ?: 2 ?? 3), "|", (true ? 2 : 3 ?? 4), "|", (false ? 1 : (false ?: 5)), "\n";
