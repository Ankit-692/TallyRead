package templates

const PasswordResetTemplate = `<!DOCTYPE html>
	<html>
	<body>
		<h2>Password Reset</h2>
		<p>Click the link below to reset your password:</p>
		<a href="%s" style="background: #007bff; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px;">Reset Password</a>
		<p>This link expires in 15 minutes.</p>
		<p>If you didn't request this, ignore this email.</p>
	</body>
	</html>
	`
