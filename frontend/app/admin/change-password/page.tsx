"use client";

import * as React from "react";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { getMe } from "@/lib/api/user";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Key, Loader2, ArrowLeft } from "lucide-react";
import { toast } from "sonner";
import api from "@/lib/api";
import Link from "next/link";

export default function ChangePasswordPage() {
	const router = useRouter();
	const [isLoading, setIsLoading] = useState(false);
	const [formData, setFormData] = useState({
		current_password: "",
		new_password: "",
		confirm_password: "",
	});

	const { data: user } = useQuery({
		queryKey: ["user-me"],
		queryFn: getMe,
	});

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();

		// Validate passwords match
		if (formData.new_password !== formData.confirm_password) {
			toast.error("New passwords do not match");
			return;
		}

		// Validate password length
		if (formData.new_password.length < 8) {
			toast.error("Password must be at least 8 characters long");
			return;
		}

		setIsLoading(true);

		try {
			await api.post("/api/v1/admin/change-password", {
				current_password: formData.current_password,
				new_password: formData.new_password,
			});

			toast.success("Password changed successfully");
			setFormData({
				current_password: "",
				new_password: "",
				confirm_password: "",
			});
		} catch (error: any) {
			const errorMsg = error.response?.data?.error || "Failed to change password";
			toast.error(errorMsg);
		} finally {
			setIsLoading(false);
		}
	};

	return (
		<div className="container mx-auto p-6 max-w-2xl font-sans min-h-screen bg-zinc-50 dark:bg-black">
			<div className="mb-8">
				<Link href="/admin">
					<Button variant="ghost" className="gap-2">
						<ArrowLeft className="h-4 w-4" />
						Back to Dashboard
					</Button>
				</Link>
			</div>

			<Card className="border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900/50 shadow-sm">
				<CardHeader>
					<div className="flex items-center gap-3">
						<div className="p-3 bg-indigo-100 dark:bg-indigo-900/40 rounded-xl">
							<Key className="h-6 w-6 text-indigo-600 dark:text-indigo-400" />
						</div>
						<div>
							<CardTitle className="text-2xl">Change Password</CardTitle>
							<CardDescription>
								{user?.email ? `Updating password for ${user.email}` : "Update your admin password"}
							</CardDescription>
						</div>
					</div>
				</CardHeader>

				<CardContent>
					<form onSubmit={handleSubmit} className="space-y-6">
						<div className="space-y-2">
							<Label htmlFor="current_password" className="text-sm font-semibold">
								Current Password
							</Label>
							<Input
								id="current_password"
								type="password"
								placeholder="Enter your current password"
								className="h-11 rounded-xl"
								value={formData.current_password}
								onChange={(e) => setFormData({ ...formData, current_password: e.target.value })}
								disabled={isLoading}
								required
							/>
						</div>

						<div className="space-y-2">
							<Label htmlFor="new_password" className="text-sm font-semibold">
								New Password
							</Label>
							<Input
								id="new_password"
								type="password"
								placeholder="Enter your new password"
								className="h-11 rounded-xl"
								value={formData.new_password}
								onChange={(e) => setFormData({ ...formData, new_password: e.target.value })}
								disabled={isLoading}
								required
								minLength={8}
							/>
							<p className="text-xs text-zinc-500">Must be at least 8 characters long</p>
						</div>

						<div className="space-y-2">
							<Label htmlFor="confirm_password" className="text-sm font-semibold">
								Confirm New Password
							</Label>
							<Input
								id="confirm_password"
								type="password"
								placeholder="Confirm your new password"
								className="h-11 rounded-xl"
								value={formData.confirm_password}
								onChange={(e) => setFormData({ ...formData, confirm_password: e.target.value })}
								disabled={isLoading}
								required
								minLength={8}
							/>
						</div>

						<Button
							type="submit"
							className="w-full h-11 rounded-xl bg-indigo-600 hover:bg-indigo-700 font-semibold"
							disabled={isLoading}
						>
							{isLoading ? (
								<>
									<Loader2 className="mr-2 h-4 w-4 animate-spin" />
									Changing Password...
								</>
							) : (
								"Change Password"
							)}
						</Button>
					</form>
				</CardContent>
			</Card>
		</div>
	);
}
