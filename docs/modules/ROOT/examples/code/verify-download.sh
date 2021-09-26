echo "🔒 Verifying download..."
{wget} \
  "{download-url}/checksums.txt" \
  "{download-url}/checksums.txt.sig" \
  "{raw-download-url}/master/.github/signature.asc"
gpg --import signature.asc
gpg --verify checksums.txt.sig
grep "$(sha256sum {artifact})" checksums.txt && echo "✔️ Checksum valid" || echo "❌ Checksum mismatch. Please investigate what went wrong."
