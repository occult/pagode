import { useState } from "react";
import AppLayout from "@/Layouts/AppLayout";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { type BreadcrumbItem } from "@/types";
import { Head, router } from "@inertiajs/react";
import { 
  CreditCardIcon, 
  CalendarIcon, 
  XCircleIcon, 
  CheckCircleIcon,
  AlertTriangleIcon 
} from "lucide-react";

interface Subscription {
  id: string;
  status: string;
  currentPeriodStart: string;
  currentPeriodEnd: string;
  amount: number;
  currency: string;
  interval: string;
  metadata?: Record<string, any>;
}

interface PaymentMethod {
  id: string;
  type: string;
  brand?: string;
  lastFour?: string;
  expiryMonth?: number;
  expiryYear?: number;
  isDefault: boolean;
}

interface User {
  id: number;
  name: string;
  email: string;
}

interface BillingProps {
  title: string;
  subscriptions: Subscription[];
  paymentMethods: PaymentMethod[];
  user: User;
}

const breadcrumbs: BreadcrumbItem[] = [
  {
    title: "Dashboard",
    href: "/dashboard",
  },
  {
    title: "Billing",
    href: "/billing",
  },
];

export default function Billing({ title, subscriptions, paymentMethods, user }: BillingProps) {
  const [isProcessing, setIsProcessing] = useState(false);
  const [processingSubscriptionId, setProcessingSubscriptionId] = useState<string | null>(null);

  const formatPrice = (price: number, currency: string) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency.toUpperCase(),
    }).format(price / 100);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  const getStatusBadge = (status: string) => {
    switch (status.toLowerCase()) {
      case 'active':
        return <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">Active</span>;
      case 'canceled':
        return <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">Canceled</span>;
      case 'past_due':
        return <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">Past Due</span>;
      case 'incomplete':
        return <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">Incomplete</span>;
      default:
        return <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">{status}</span>;
    }
  };

  const handleCancelSubscription = (subscriptionId: string) => {
    if (!confirm('Are you sure you want to cancel this subscription? This action cannot be undone.')) {
      return;
    }

    setIsProcessing(true);
    setProcessingSubscriptionId(subscriptionId);
    
    router.post('/billing/cancel', {
      subscriptionId: subscriptionId,
    }, {
      onSuccess: () => {
        // The backend redirects, so we don't need to handle success here
      },
      onError: (errors) => {
        console.error('Cancellation failed:', errors);
        alert('Failed to cancel subscription. Please try again.');
      },
      onFinish: () => {
        setIsProcessing(false);
        setProcessingSubscriptionId(null);
      },
    });
  };

  const hasActiveSubscriptions = subscriptions.some(sub => sub.status === 'active');

  return (
    <AppLayout breadcrumbs={breadcrumbs}>
      <Head title={title} />
      <div className="flex h-full flex-1 flex-col gap-6 rounded-xl p-6">
        {/* Header */}
        <div className="space-y-2">
          <h1 className="text-3xl font-bold tracking-tight">Billing & Subscription</h1>
          <p className="text-muted-foreground">
            Manage your subscription and payment methods
          </p>
        </div>

        {/* No Subscriptions Message */}
        {subscriptions.length === 0 && (
          <div className="rounded-lg border border-yellow-200 bg-yellow-50 p-4">
            <div className="flex items-center gap-3">
              <AlertTriangleIcon className="h-4 w-4 text-yellow-600" />
              <div>
                <p className="text-sm text-yellow-800">
                  You don't have any active subscriptions. 
                  <Button variant="link" className="p-0 ml-2 h-auto" asChild>
                    <a href="/plans">View available plans</a>
                  </Button>
                </p>
              </div>
            </div>
          </div>
        )}

        {/* Active Subscriptions */}
        {subscriptions.length > 0 && (
          <div className="space-y-4">
            <h2 className="text-xl font-semibold">Your Subscriptions</h2>
            <div className="grid gap-4">
              {subscriptions.map((subscription) => (
                <Card key={subscription.id}>
                  <CardHeader>
                    <div className="flex items-center justify-between">
                      <div className="space-y-1">
                        <CardTitle className="flex items-center gap-2">
                          <CreditCardIcon className="h-5 w-5" />
                          {subscription.metadata?.plan_id || 'Premium Plan'}
                          {getStatusBadge(subscription.status)}
                        </CardTitle>
                        <CardDescription>
                          {formatPrice(subscription.amount, subscription.currency)}/{subscription.interval}
                        </CardDescription>
                      </div>
                      {subscription.status === 'active' && (
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleCancelSubscription(subscription.id)}
                          disabled={isProcessing && processingSubscriptionId === subscription.id}
                        >
                          {isProcessing && processingSubscriptionId === subscription.id ? (
                            <div className="flex items-center gap-2">
                              <div className="h-3 w-3 animate-spin rounded-full border-2 border-current border-t-transparent" />
                              Canceling...
                            </div>
                          ) : (
                            <>
                              <XCircleIcon className="h-4 w-4 mr-2" />
                              Cancel
                            </>
                          )}
                        </Button>
                      )}
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="space-y-2">
                        <div className="flex items-center gap-2 text-sm">
                          <CalendarIcon className="h-4 w-4 text-muted-foreground" />
                          <span className="text-muted-foreground">Current Period:</span>
                        </div>
                        <p className="text-sm font-medium">
                          {formatDate(subscription.currentPeriodStart)} - {formatDate(subscription.currentPeriodEnd)}
                        </p>
                      </div>
                      <div className="space-y-2">
                        <div className="flex items-center gap-2 text-sm">
                          <CheckCircleIcon className="h-4 w-4 text-muted-foreground" />
                          <span className="text-muted-foreground">Status:</span>
                        </div>
                        <p className="text-sm font-medium capitalize">{subscription.status}</p>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>
        )}

        {/* Payment Methods */}
        {paymentMethods.length > 0 && (
          <div className="space-y-4">
            <h2 className="text-xl font-semibold">Payment Methods</h2>
            <div className="grid gap-4">
              {paymentMethods.map((method) => (
                <Card key={method.id}>
                  <CardContent className="pt-6">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-3">
                        <CreditCardIcon className="h-5 w-5 text-muted-foreground" />
                        <div>
                          <p className="font-medium">
                            {method.brand?.toUpperCase() || method.type.toUpperCase()} ••••{method.lastFour}
                          </p>
                          {method.expiryMonth && method.expiryYear && (
                            <p className="text-sm text-muted-foreground">
                              Expires {method.expiryMonth.toString().padStart(2, '0')}/{method.expiryYear}
                            </p>
                          )}
                        </div>
                      </div>
                      {method.isDefault && (
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">Default</span>
                      )}
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>
        )}

        {/* Quick Actions */}
        <Card>
          <CardHeader>
            <CardTitle>Quick Actions</CardTitle>
            <CardDescription>
              Need to make changes to your subscription or billing?
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-3">
              {!hasActiveSubscriptions && (
                <Button asChild>
                  <a href="/plans">
                    <CreditCardIcon className="h-4 w-4 mr-2" />
                    Subscribe to a Plan
                  </a>
                </Button>
              )}
              <Button variant="outline" asChild>
                <a href="/plans">
                  <CalendarIcon className="h-4 w-4 mr-2" />
                  View All Plans
                </a>
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Account Information */}
        <Card>
          <CardHeader>
            <CardTitle>Account Information</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Name:</span>
                <span className="font-medium">{user.name}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Email:</span>
                <span className="font-medium">{user.email}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Customer ID:</span>
                <span className="font-medium">#{user.id}</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </AppLayout>
  );
}