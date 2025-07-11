import { useState } from "react";
import AppLayout from "@/Layouts/AppLayout";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { PaymentForm } from "@/components/PaymentForm";
import { type BreadcrumbItem } from "@/types";
import { Head, router } from "@inertiajs/react";
import { CheckIcon, CreditCardIcon } from "lucide-react";

interface Plan {
  id: string;
  name: string;
  description: string;
  price: number;
  currency: string;
  interval: string;
  features: string[];
  priceId: string;
}

interface Form {
  data: {
    planId: string;
    paymentMethodId: string;
  };
  errors: Record<string, string>;
  isSubmitted: boolean;
  isValid: boolean;
}

interface PlansProps {
  title: string;
  hasActiveSubscription: boolean;
  plans: Plan[];
  form: Form;
  stripePublishableKey: string;
}

const breadcrumbs: BreadcrumbItem[] = [
  {
    title: "Dashboard",
    href: "/dashboard",
  },
  {
    title: "Plans",
    href: "/plans",
  },
];

export default function Plans({ title, hasActiveSubscription, plans, form, stripePublishableKey }: PlansProps) {
  const [selectedPlan, setSelectedPlan] = useState<Plan | null>(null);
  const [isProcessing, setIsProcessing] = useState(false);
  const [showPaymentModal, setShowPaymentModal] = useState(false);

  const formatPrice = (price: number, currency: string) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency.toUpperCase(),
    }).format(price / 100);
  };

  const handleSelectPlan = (plan: Plan) => {
    setSelectedPlan(plan);
    setShowPaymentModal(true);
  };

  const handlePaymentSubmit = (paymentMethodId: string) => {
    if (!selectedPlan) return;
    
    setIsProcessing(true);
    
    router.post('/plans/subscribe', {
      planId: selectedPlan.id,
      paymentMethodId: paymentMethodId,
    }, {
      onSuccess: () => {
        setShowPaymentModal(false);
        setSelectedPlan(null);
      },
      onError: (errors) => {
        console.error('Subscription failed:', errors);
        // Form errors are now handled by the backend and redisplayed
      },
      onFinish: () => {
        setIsProcessing(false);
      },
    });
  };

  return (
    <AppLayout breadcrumbs={breadcrumbs}>
      <Head title={title} />
      <div className="flex h-full flex-1 flex-col gap-6 rounded-xl p-6">
        {/* Header */}
        <div className="text-center space-y-2">
          <h1 className="text-3xl font-bold tracking-tight">Choose Your Plan</h1>
          <p className="text-muted-foreground">
            Unlock premium features and take your experience to the next level
          </p>
        </div>

        {/* Active Subscription Notice */}
        {hasActiveSubscription && (
          <Card className="border-green-200 bg-green-50 dark:border-green-800 dark:bg-green-950">
            <CardContent className="flex items-center gap-3 pt-6">
              <CheckIcon className="h-5 w-5 text-green-600" />
              <div>
                <p className="font-medium text-green-800 dark:text-green-200">
                  You have an active subscription
                </p>
                <p className="text-sm text-green-600 dark:text-green-400">
                  Thank you for being a premium member!
                </p>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Plans Grid */}
        <div className="grid gap-6 md:grid-cols-1 lg:grid-cols-3 max-w-4xl mx-auto w-full">
          {plans.map((plan) => (
            <Card
              key={plan.id}
              className={`relative transition-all duration-200 hover:shadow-lg ${
                selectedPlan?.id === plan.id ? 'ring-2 ring-primary' : ''
              }`}
            >
              <CardHeader className="text-center">
                <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
                  <CreditCardIcon className="h-6 w-6 text-primary" />
                </div>
                <CardTitle className="text-2xl">{plan.name}</CardTitle>
                <CardDescription>{plan.description}</CardDescription>
                <div className="mt-4">
                  <span className="text-4xl font-bold">
                    {formatPrice(plan.price, plan.currency)}
                  </span>
                  <span className="text-muted-foreground">/{plan.interval}</span>
                </div>
              </CardHeader>
              
              <CardContent className="space-y-4">
                {/* Features List */}
                <ul className="space-y-3">
                  {plan.features.map((feature, index) => (
                    <li key={index} className="flex items-center gap-3">
                      <CheckIcon className="h-4 w-4 text-green-600 flex-shrink-0" />
                      <span className="text-sm">{feature}</span>
                    </li>
                  ))}
                </ul>

                {/* Subscribe Button */}
                <Button
                  className="w-full mt-6"
                  onClick={() => handleSelectPlan(plan)}
                  disabled={isProcessing || hasActiveSubscription}
                  variant={hasActiveSubscription ? "outline" : "default"}
                >
                  {isProcessing && selectedPlan?.id === plan.id ? (
                    <div className="flex items-center gap-2">
                      <div className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                      Processing...
                    </div>
                  ) : hasActiveSubscription ? (
                    "Already Subscribed"
                  ) : (
                    `Subscribe to ${plan.name}`
                  )}
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Form Errors */}
        {form.isSubmitted && !form.isValid && (
          <Card className="max-w-2xl mx-auto border-red-200 bg-red-50 dark:border-red-800 dark:bg-red-950">
            <CardContent className="pt-6">
              <div className="flex items-start gap-3">
                <XCircleIcon className="h-5 w-5 text-red-600 flex-shrink-0 mt-0.5" />
                <div className="space-y-1">
                  <p className="font-medium text-red-800 dark:text-red-200">
                    Subscription Failed
                  </p>
                  <div className="text-sm text-red-700 dark:text-red-300">
                    {Object.entries(form.errors).map(([field, error]) => (
                      <p key={field}>{error}</p>
                    ))}
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        )}

        {/* Demo Notice */}
        <Card className="max-w-2xl mx-auto border-yellow-200 bg-yellow-50 dark:border-yellow-800 dark:bg-yellow-950">
          <CardContent className="pt-6">
            <div className="flex items-start gap-3">
              <div className="h-5 w-5 rounded-full bg-yellow-400 flex-shrink-0 mt-0.5" />
              <div className="space-y-1">
                <p className="font-medium text-yellow-800 dark:text-yellow-200">
                  Demo Mode
                </p>
                <p className="text-sm text-yellow-700 dark:text-yellow-300">
                  This is a demonstration of the payment system. In a real implementation, 
                  you would integrate Stripe Elements for secure payment method collection.
                  The subscription endpoint is functional and ready for production use.
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Payment Modal */}
        <Dialog open={showPaymentModal} onOpenChange={setShowPaymentModal}>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>Complete Your Subscription</DialogTitle>
            </DialogHeader>
            {selectedPlan && (
              <PaymentForm
                plan={selectedPlan}
                onSubmit={handlePaymentSubmit}
                isProcessing={isProcessing}
                stripePublishableKey={stripePublishableKey}
              />
            )}
          </DialogContent>
        </Dialog>
      </div>
    </AppLayout>
  );
}